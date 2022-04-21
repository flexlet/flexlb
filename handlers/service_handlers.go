// Copyright (c) 2022 Yaohui Wang (yaohuiwang@outlook.com)
// FlexLB is licensed under Mulan PubL v2.
// You can use this software according to the terms and conditions of the Mulan PubL v2.
// You may obtain a copy of Mulan PubL v2 at:
//         http://license.coscl.org.cn/MulanPubL-2.0
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PubL v2 for more details.

package handlers

import (
	"gitee.com/flexlb/flexlb-api/config"
	"gitee.com/flexlb/flexlb-api/restapi/operations/service"

	"github.com/go-openapi/runtime/middleware"
)

type ReadyzHandlerImpl struct {
}

func (h *ReadyzHandlerImpl) Handle(params service.ReadyzParams) middleware.Responder {
	return service.NewReadyzOK().WithPayload(config.LB.Status)
}
