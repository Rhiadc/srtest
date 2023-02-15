package resources

type CustomResourceDefinition struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`

	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`

	Spec struct {
		Image        string `yaml:"image"`
		FlinkVersion string `yaml:"flinkVersion"`

		FlinkConfiguration struct {
			TaskmanagerNumberOfTaskSlots string `json:"taskmanager.numberOfTaskSlots"`
			StateSavepointsDir           string `json:"state.savepoints.dir"`
			StateCheckpointsDir          string `json:"state.checkpoints.dir"`
		} `yaml:"flinkConfiguration"`

		ServiceAccount string `yaml:"serviceAccount"`

		PodTemplate struct {
			APIVersion string `yaml:"apiVersion"`
			Kind       string `yaml:"kind"`

			Metadata struct {
				Name string `yaml:"name"`
			} `yaml:"metadata"`

			Spec struct {
				ServiceAccount string `yaml:"serviceAccount"`
				Tolerations    []struct {
					Effect   string `yaml:"effect"`
					Key      string `yaml:"key"`
					Operator string `yaml:"operator"`
					Value    string `yaml:"value"`
				} `yaml:"tolerations"`
			} `yaml:"spec"`
		} `yaml:"podTemplate"`

		JobManager struct {
			Resource struct {
				Memory string `yaml:"memory"`
				CPU    int    `yaml:"cpu"`
			} `yaml:"resource"`

			PodTemplate struct {
				APIVersion string `yaml:"apiVersion"`
				Kind       string `yaml:"kind"`
				Spec       struct {
					Affinity struct {
						NodeAffinity struct {
							RequiredDuringSchedulingIgnoredDuringExecution struct {
								NodeSelectorTerms []struct {
									MatchExpressions []struct {
										Key      string   `yaml:"key"`
										Operator string   `yaml:"operator"`
										Values   []string `yaml:"values"`
									} `yaml:"matchExpressions"`
								} `yaml:"nodeSelectorTerms"`
							} `yaml:"requiredDuringSchedulingIgnoredDuringExecution"`
						} `yaml:"nodeAffinity"`
					} `yaml:"affinity"`
				} `yaml:"spec"`
			} `yaml:"podTemplate"`
		} `yaml:"jobManager"`

		TaskManager struct {
			Resource struct {
				Memory string `yaml:"memory"`
				CPU    int    `yaml:"cpu"`
			} `yaml:"resource"`
		} `yaml:"taskManager"`

		Job struct {
			JarURI      string   `yaml:"jarURI"`
			Parallelism int      `yaml:"parallelism"`
			UpgradeMode string   `yaml:"upgradeMode"`
			State       string   `yaml:"state"`
			Args        []string `yaml:"args"`
		} `yaml:"job"`
	} `yaml:"spec"`
}
