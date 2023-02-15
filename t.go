package main

import (
	"log"

	"github.com/rhiadc/srtest/git"
	"github.com/rhiadc/srtest/yaml"
)

var jsonData = []byte(`
{
	"apiVersion": "flink.apache.org/v1beta1",
	"kind": "FlinkDeployment",
	"metadata": {
	  "name": "my-job-name"
	},
	"spec": {
	  "image": "my-image-name:eee1.0.0",
	  "flinkVersion": "v1_14",
	  "flinkConfiguration": {
		"taskmanager.numberOfTaskSlots": "2",
		"state.savepoints.dir": "ointing-qa/my-job-name/",
		"state.checkpoints.dir": "nk-checkpointing-qa/my-job-name/savepoint"
	  },
	  "serviceAccount": "flink",
	  "podTemplate": {
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": {
		  "name": "pod-template"
		},
		"spec": {
		  "serviceAccount": "flink",
		  "tolerations": [
			{
			  "effect": "NoExecute",
			  "key": "spot",
			  "operator": "Equal",
			  "value": "true"
			},
			{
			  "effect": "NoSchedule",
			  "key": "spot",
			  "operator": "Equal",
			  "value": "true"
			}
		  ]
		}
	  },
	  "jobManager": {
		"resource": {
		  "memory": "2048m",
		  "cpu": 1
		},
		"podTemplate": {
		  "apiVersion": "v1",
		  "kind": "Pod",
		  "spec": {
			"affinity": {
			  "nodeAffinity": {
				"requiredDuringSchedulingIgnoredDuringExecution": {
				  "nodeSelectorTerms": [
					{
					  "matchExpressions": [
						{
						  "key": "family",
						  "operator": "In",
						  "values": [
							"value"
						  ]
						},
						{
						  "key": "team-ops/workload",
						  "operator": "In",
						  "values": [
							"value"
						  ]
						},
						{
						  "key": "app",
						  "operator": "In",
						  "values": [
							"value"
						  ]
						}
					  ]
					}
				  ]
				}
			  }
			}
		  }
		}
	  },
	  "taskManager": {
		"resource": {
		  "memory": "2048m",
		  "cpu": 1
		}
	  },
	  "job": {
		"jarURI": "local:///opt/flink/usrlib/my-flink-job.jar",
		"parallelism": 2,
		"upgradeMode": "savepoint",
		"state": "running",
		"args": [
		  "--source-servers",
		  "b-5.k",
		  "--source-topic",
		  "complex",
		  "--source-group-id",
		  "my-job-name-op",
		  "--sink-topic",
		  "my-job-name",
		  "--sink-servers",
		  "joba,jobc,joba"
		]
	  }
	}
  }`)

func main() {
	gconfig := git.GitConfig{
		RepoURL:           "https://github.com/Rhiadc/gitops-testapi.git",
		RepoInternalPath:  "pp",
		RepoReferenceName: "refs/heads/main",
		Username:          "Rhiadc",
		Password:          "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",
		OwnerName:         "rhiad",
		OnwerEmail:        "rhiad.ciccoli@gmail.com",
	}

	gitClient := git.NewGitClient(gconfig)
	repo, err := gitClient.CloneRepo()
	if err != nil {
		log.Fatal(err)
	}

	wt, err := gitClient.PullFromMain(repo)
	if err != nil {
		log.Fatal(err)
	}

	fileName, err := yaml.GenerateK8SYAMLAndValidate(jsonData)
	log.Fatal(err)

	if err := gitClient.AddCommitAndPush(repo, wt, fileName); err != nil {
		log.Fatal(err)
	}

}
