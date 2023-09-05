package routes

import "github.com/yonyoucloud/datatable/pkg/modules"

type Routes struct {
	modules *modules.Modules
}

func New(ms *modules.Modules) *Routes {
	return &Routes{modules: ms}
}
