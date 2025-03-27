package bo

type AppsInfo struct {
	Items []struct {
		Metadata *Metadata `json:"metadata"`
		Spec     *Spec     `json:"spec"`
		Status   struct {
			Health         *Health         `json:"health"`
			OperationState *OperationState `json:"operationState"`
			Summary        *Summary        `json:"summary"`
			Sync           *Sync           `json:"sync"`
		} `json:"status"`
	} `json:"items"`
}

type AppInfo struct {
	Metadata *Metadata `json:"metadata"`
	Status   struct {
		History        []*History      `json:"history"`
		OperationState *OperationState `json:"operationState"`
		Summary        *Summary        `json:"summary"`
	} `json:"status"`
}

type Metadata struct {
	CreationTimestamp string `json:"creationTimestamp"`
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
}

type Spec struct {
	Destination struct {
		Namespace string `json:"namespace"`
		Server    string `json:"server"`
	} `json:"destination"`
	Project              string `json:"project"`
	RevisionHistoryLimit string `json:"revisionHistoryLimit"`
	Source               struct {
		Helm struct {
			ReleaseName string   `json:"releaseName"`
			ValueFiles  []string `json:"valueFiles"`
		} `json:"helm"`
		Path           string `json:"path"`
		RepoURL        string `json:"repoURL"`
		TargetRevision string `json:"targetRevision"`
	} `json:"source"`
}

type Health struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type History struct {
	DeployStartedAt string `json:"deployStartedAt"`
	DeployedAt      string `json:"deployedAt"`
	Id              int    `json:"id"`
	Revision        string `json:"revision"`
	Source          struct {
		Helm struct {
			ReleaseName string   `json:"releaseName"`
			ValueFiles  []string `json:"valueFiles"`
		} `json:"helm"`
		Path           string `json:"path"`
		RepoURL        string `json:"repoURL"`
		TargetRevision string `json:"targetRevision"`
	} `json:"source"`
}

type OperationState struct {
	StartedAt  string `json:"startedAt"`
	FinishedAt string `json:"finishedAt"`
}

type Summary struct {
	ExternalURLs []string `json:"externalURLs"`
	Images       []string `json:"images"`
}

type Image struct {
	Helm struct {
		Parameters []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"parameters"`
	} `json:"helm"`
	Type string `json:"type"`
}

type ImageInfo struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
}

type Sync struct {
	Revision string `json:"revision"`
	Status   string `json:"status"`
}

type Log struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
	Result *Result `json:"result"`
}

type Result struct {
	Content string `json:"content"`
}

type ResourceTree struct {
	Nodes []struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
	} `json:"nodes"`
}
