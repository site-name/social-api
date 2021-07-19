package plugin

import (
	"github.com/pkg/errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

// CheckRequiredServerConfiguration implements Helpers.CheckRequiredServerConfiguration
func (p *HelpersImpl) CheckRequiredServerConfiguration(req *model.Config) (bool, error) {
	if req == nil {
		return true, nil
	}

	cfg := p.API.GetConfig()

	mc, err := util.Merge(cfg, req, nil)
	if err != nil {
		return false, errors.Wrap(err, "could not merge configurations")
	}

	mergedCfg := mc.(model.Config)
	if mergedCfg.ToJson() != cfg.ToJson() {
		return false, nil
	}

	return true, nil
}
