package com.vincent.service;

import com.vincent.pojo.DbStu;

import java.util.List;

public interface StuService {
    /**
     * 新增stu到数据库
     * @param stu
     */
    public void saveStu(DbStu stu);

    /**
     * 根据主键查询对象信息
     * @param id
     * @return DbStu
     */
    public DbStu queryById(String id);

    /**
     * 根据条件查询stu的list结果集
     * @param name
     * @param sex
     * @return List<DbStu>
     */
    public List<DbStu> queryByCondition(String name, Integer sex);
}
