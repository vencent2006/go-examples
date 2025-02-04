package com.vincent.huobi.client.req.subuser;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import com.vincent.huobi.constant.enums.SubUserManagementActionEnum;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class SubUserManagementRequest {

  private Long subUid;

  private SubUserManagementActionEnum action;

}
