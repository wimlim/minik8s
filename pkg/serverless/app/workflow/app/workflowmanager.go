package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/message"
	"net/http"
	"strings"
	"time"

	"github.com/knetic/govaluate"
	"github.com/streadway/amqp"
)

const urlbase = "http://localhost:8081"

func StartWorkflow(w *apiobj.Workflow) {
	nodebase := w.Spec.WorkflowNodes
	entryname := w.Spec.EntryName
	fmt.Println("Start workflow")
	curNode := findNodeByName(nodebase, entryname)
	if curNode == nil {
		fmt.Println("Entry node not found")
		return
	}
	curParam := findParam(curNode.FuncNode.FuncParam, w.Spec.EntryParam)

	var result map[string]interface{}
	var flag bool
	var err error
	for curNode != nil {
		switch curNode.Type {
		case apiobj.FunctionType:
			for {
				result, err = execFuncNode(curNode, curParam)
				if err == nil {
					break
				}
				time.Sleep(3 * time.Second)
			}
			if curNode.FuncNode.NextNodeName == apiobj.EndNode {
				fmt.Println("End of workflow")
				return
			}
			curNode = findNodeByName(nodebase, curNode.FuncNode.NextNodeName)
			if curNode.Type == apiobj.FunctionType {
				curParam = findParam(curNode.FuncNode.FuncParam, result)
			}
		case apiobj.ChoiceType:
			flag, err = execChoiceNode(curNode, result)
			if err != nil {
				fmt.Println(err)
				return
			}
			if flag {
				if curNode.ChoiceNode.TrueNodeName == apiobj.EndNode {
					fmt.Println("End of workflow")
					return
				}
				curParam = fillParam(curNode.ChoiceNode.TrueEntryParam, result)
				curNode = findNodeByName(nodebase, curNode.ChoiceNode.TrueNodeName)
			} else {
				if curNode.ChoiceNode.FalseNodeName == apiobj.EndNode {
					fmt.Println("End of workflow")
					return
				}
				curParam = fillParam(curNode.ChoiceNode.FalseEntryParam, result)
				curNode = findNodeByName(nodebase, curNode.ChoiceNode.FalseNodeName)
			}
		}
	}
}

func findNodeByName(nodes []apiobj.WorkflowNode, name string) *apiobj.WorkflowNode {
	for _, node := range nodes {
		if node.Name == name {
			return &node
		}
	}
	return nil
}

func findParam(params []string, input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, param := range params {
		result[param] = input[param]
	}
	return result
}

func fillParam(params map[string]interface{}, input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	// iterate over the params and fill the values
	for key, value := range params {
		if valueStr, ok := value.(string); ok {
			if valueStr[0] == '$' {
				result[key] = input[valueStr[1:]]
			} else {
				result[key] = value
			}
		} else {
			result[key] = value
		}
	}
	return result
}

func execFuncNode(node *apiobj.WorkflowNode, param map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println("Execute ", node.Name)
	url := fmt.Sprintf("%s/%s/%s", urlbase, node.FuncNode.FuncNameSpace, node.FuncNode.FuncName)
	requestBody, err := json.Marshal(param)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Function invocation failed")
		return nil, fmt.Errorf("function invocation failed: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	return result, nil
}

func execChoiceNode(node *apiobj.WorkflowNode, param map[string]interface{}) (bool, error) {
	fmt.Println("Execute ", node.Name)
	exprStr := node.ChoiceNode.Expression
	// replace the expression with actual values
	fmt.Println(node)
	fmt.Println(param)
	for key, value := range param {
		if valueStr, ok := value.(string); ok {
			exprStr = strings.Replace(exprStr, "$"+key, valueStr, -1)
		} else {
			exprStr = strings.Replace(exprStr, "$"+key, fmt.Sprintf("%v", value), -1)
		}
	}
	fmt.Println("Expression: ", exprStr)
	// execute the expression
	expr, err := govaluate.NewEvaluableExpression(exprStr)
	if err != nil {
		return false, err
	}
	result, err := expr.Evaluate(param)
	if err != nil {
		return false, err
	}
	if boolResult, ok := result.(bool); ok {
		return boolResult, nil
	}
	return false, nil
}

func Run() {
	sub := message.NewSubscriber()
	defer sub.Close()
	sub.Subscribe(message.WorkflowQueue, func(d amqp.Delivery) {
		var message message.Message
		json.Unmarshal([]byte(d.Body), &message)

		var workflow apiobj.Workflow
		json.Unmarshal([]byte(message.Content), &workflow)

		StartWorkflow(&workflow)
	})
}
