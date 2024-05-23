package service

import "github.com/google/wire"

var InjectSet = wire.NewSet(NewTagService, NewCategoryService, NewConcatService, NewLinkService, NewMenuService, NewSysConfigService,
	NewOtherService)
