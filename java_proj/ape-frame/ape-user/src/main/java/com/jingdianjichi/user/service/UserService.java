package com.jingdianjichi.user.service;

import com.jingdianjichi.entity.PageResult;
import com.jingdianjichi.user.entity.dto.UserDto;
import com.jingdianjichi.user.entity.po.UserPo;

// controller -> service 用dto
// service -> mapper 用entity(po)
public interface UserService {
    int addUser(UserDto userDto);

    int delete(Integer id);

    PageResult<UserPo> getUserPage(UserDto userDto);
}
