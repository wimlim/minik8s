package jobserver

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/jobserverconfig"
	"minik8s/pkg/message"
	"os"
	"os/exec"

	"github.com/streadway/amqp"
	"golang.org/x/crypto/ssh"
)

const sshlocation = "stu1938@pilogin.hpc.sjtu.edu.cn:/lustre/home/acct-stu/stu1938/"
const slurmlocation = "/tmp/jobs/"

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
	location := sshlocation + "job-" + job.MetaData.Name
	// scp file to server
	cmd := exec.Command("scp", "-r", job.File, location)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to scp file: ", err)
		return
	}
	// create slurm script in /tmp/jobs/name.slurm
	script := generateSlurmScript(job)
	err = os.WriteFile(slurmlocation+job.MetaData.Name+".slurm", []byte(script), 0644)
	if err != nil {
		fmt.Println("Failed to write slurm script: ", err)
		return
	}
	// scp slurm
	cmd = exec.Command("scp", "-r", slurmlocation+job.MetaData.Name+".slurm", location)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Failed to scp slurm script: ", err)
		return
	}
	// run slurm
	session, err := js.sshClient.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: ", err)
		return
	}
	defer session.Close()
	err = session.Run("sbatch " + location + "/" + job.MetaData.Name + ".slurm")
	if err != nil {
		fmt.Println("Failed to run slurm script: ", err)
		return
	}
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
	nvcc ` + job.File + ` -o ` + job.MetaData.Name + `
	./` + job.MetaData.Name
	return script
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
			jobServer.CreateJob(&job)
		case "Query":

		}

	})
}
