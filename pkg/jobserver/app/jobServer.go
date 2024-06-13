package jobserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"minik8s/pkg/config/jobserverconfig"
	"minik8s/pkg/message"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"golang.org/x/crypto/ssh"
)

const sshlocation = "stu1938@pilogin.hpc.sjtu.edu.cn:/lustre/home/acct-stu/stu1938/"
const slurmlocation = "/tmp/jobs/"
const reslocation = "/tmp/results/"

type JobServer struct {
	sshClient  *ssh.Client
	subscriber *message.Subscriber
}

func NewJobServer() *JobServer {
	sshConfig := &ssh.ClientConfig{
		User: jobserverconfig.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(jobserverconfig.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", jobserverconfig.HttpUrl+":22", sshConfig)
	if err != nil {
		fmt.Println("Failed to dial: ", err)
	}
	subscriber := message.NewSubscriber()
	return &JobServer{
		sshClient:  sshClient,
		subscriber: subscriber,
	}
}

func (js *JobServer) CreateJob(job *apiobj.Job) {
	workspace := sshlocation + job.MetaData.Name
	// mkdir
	if err := js.runCommand("mkdir -p " + job.MetaData.Name); err != nil {
		fmt.Println("Failed to create workspace: ", err)
		return
	}
	// scp file to server
	if err := scpFile(job.File, workspace); err != nil {
		fmt.Println("Failed to scp file: ", err)
		return
	}
	// create slurm script in /tmp/jobs/name.slurm
	scriptPath := slurmlocation + job.MetaData.Name + ".slurm"
	script := generateSlurmScript(job)
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		fmt.Println("Failed to write slurm script: ", err)
		return
	}
	// scp slurm
	if err := scpFile(scriptPath, workspace); err != nil {
		fmt.Println("Failed to scp slurm script: ", err)
		return
	}
	// run slurm
	slurmCommand := fmt.Sprintf("cd %s && sbatch %s.slurm", job.MetaData.Name, job.MetaData.Name)
	if err := js.runCommand(slurmCommand); err != nil {
		fmt.Println("Failed to run slurm script: ", err)
		return
	}
}

func (js *JobServer) MonitorJob(job *apiobj.Job) bool {
	timeout := time.After(1 * time.Hour)
	ticker := time.NewTicker(5 * time.Second)
	tick := ticker.C
	defer ticker.Stop()
	dir := sshlocation + job.MetaData.Name + "/"
	state := apiobj.Pending
	for {
		select {
		case <-timeout:
			fmt.Println("Job timeout")
			return false
		case <-tick:
			switch js.getJobStatus(job.MetaData.Name) {
			case apiobj.Pending:
				if state != apiobj.Pending {
					state = apiobj.Pending
					job.Status.Phase = apiobj.Pending
					fmt.Println("Job pending")
					js.updateJobStatus(job)
				}
			case apiobj.Running:
				if state != apiobj.Running {
					fmt.Println("Job running")
					state = apiobj.Running
					job.Status.Phase = apiobj.Running
					js.updateJobStatus(job)
				}
			case apiobj.Finished:
				fmt.Println("Job finished")
				// scp name.err & name.out
				if err := scpFile(dir+job.MetaData.Name+".err", reslocation+job.MetaData.Name+".err"); err != nil {
					fmt.Println("Failed to scp err file: ", err)
				}
				if err := scpFile(dir+job.MetaData.Name+".out", reslocation+job.MetaData.Name+".out"); err != nil {
					fmt.Println("Failed to scp out file: ", err)
				}
				// rmdir
				if err := js.runCommand("rm -rf " + job.MetaData.Name); err != nil {
					fmt.Println("Failed to remove workspace: ", err)
				}
				job.Status.Phase = apiobj.Finished
				js.updateJobStatus(job)
				return true
			}
		}
	}
}

func (js *JobServer) getJobStatus(name string) string {
	command := fmt.Sprintf("squeue | grep %s", name)
	session, err := js.sshClient.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: ", err)
		return "error"
	}
	defer session.Close()
	output, err := session.Output(command)
	if err != nil {
		return apiobj.Finished
	}
	if strings.Contains(string(output), "PD") {
		return apiobj.Pending
	}
	if strings.Contains(string(output), "R") {
		return apiobj.Running
	}
	return apiobj.Finished
}

func (js *JobServer) runCommand(command string) error {
	session, err := js.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()
	if err := session.Run(command); err != nil {
		return fmt.Errorf("command failed: %v", err)
	}
	return nil
}

func scpFile(localPath, remotePath string) error {
	cmd := exec.Command("scp", "-r", localPath, remotePath)
	return cmd.Run()
}

func generateSlurmScript(job *apiobj.Job) string {
	script := `#!/bin/bash
#SBATCH --job-name=` + job.MetaData.Name + `
#SBATCH --partition=` + job.Spec.Partition + `
#SBATCH --output=` + job.MetaData.Name + `.out
#SBATCH --error=` + job.MetaData.Name + `.err
#SBATCH -N ` + fmt.Sprint(job.Spec.Nodes) + `
#SBATCH --ntasks-per-node=` + fmt.Sprint(job.Spec.NtasksPerNode) + `
#SBATCH --cpus-per-task=` + fmt.Sprint(job.Spec.CpusPerTask) + `
#SBATCH --gres=` + job.Spec.Gres + `
ulimit -s unlimited
ulimit -l unlimited

module load cuda/12.1.1
` + job.Script
	return script
}

func (js *JobServer) updateJobStatus(job *apiobj.Job) {
	URL := apiconfig.URL_JobStatus
	URL = strings.Replace(URL, ":namespace", job.MetaData.Namespace, -1)
	URL = strings.Replace(URL, ":name", job.MetaData.Name, -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL
	jsonData, err := json.Marshal(job.Status)
	if err != nil {
		fmt.Println("marshal job error")
		return
	}
	req, err := http.NewRequest(http.MethodPut, HttpUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("create put request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("put error:", err)
		return
	}
	defer response.Body.Close()
}

func Run() {
	jobServer := NewJobServer()
	defer jobServer.sshClient.Close()
	defer jobServer.subscriber.Close()

	fmt.Println("JobServer is running")

	jobServer.subscriber.Subscribe(message.JobQueue, func(d amqp.Delivery) {
		var message message.Message
		err := json.Unmarshal(d.Body, &message)
		if err != nil {
			fmt.Println("Failed to unmarshal message: ", err)
		}
		var job apiobj.Job
		if err := json.Unmarshal([]byte(message.Content), &job); err != nil {
			fmt.Println("Failed to unmarshal job: ", err)
		}
		switch message.Type {
		case "Add":
			go func(job apiobj.Job) {
				jobServer.CreateJob(&job)
				if !jobServer.MonitorJob(&job) {
					fmt.Println("Job failed")
				}
			}(job)
		}
	})
}
