package com.jingdianjichi.user.entity.po;

import lombok.Data;

import java.io.Serializable;
import java.util.Date;

/**
 * (SysUser)实体类
 *
 * @author makejava
 * @since 2025-01-11 00:35:38
 */
@Data
public class SysUser implements Serializable {

    private Long id;

    private String name;

    private Integer age;

    private String createBy;

    private Date createTime;

    private String updateBy;

    private Date updateTime;

    private Integer deleteFlag;

    private Integer version;
}

