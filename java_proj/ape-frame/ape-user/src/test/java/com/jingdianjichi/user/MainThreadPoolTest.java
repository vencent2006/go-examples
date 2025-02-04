package com.jingdianjichi.user;

import com.alibaba.fastjson.JSON;
import com.jingdianjichi.redis.util.RedisShareLockUtil;
import com.jingdianjichi.tool.CompletableFutureUtils;
import com.jingdianjichi.user.delayQueue.MassMailTask;
import com.jingdianjichi.user.delayQueue.MassMailTaskService;
import com.jingdianjichi.user.event.Person;
import com.jingdianjichi.user.event.PersonEventService;
import lombok.extern.slf4j.Slf4j;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import javax.annotation.Resource;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.UUID;
import java.util.concurrent.FutureTask;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

// classes 表示启动类
// webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT 表示随机端口启动应用
@SpringBootTest(classes = UserApplication.class, webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@RunWith(SpringRunner.class) // 表示启动这个单元测试类（单元测试类是不能够运行的），需要传递一个参数，必须是SpringRunner的实例类型
@Slf4j
public class MainThreadPoolTest {
    @Resource(name = "mailThreadPool")
    private ThreadPoolExecutor mailThreadPool;

    @Resource
    private PersonEventService personEventService;

    @Resource
    private MassMailTaskService massMailTaskService;

    @Resource
    private RedisShareLockUtil redisShareLockUtil;

    @Test
    public void test() {
        for (int i = 0; i < 10; i++) {
            mailThreadPool.submit(new Runnable() {
                @Override
                public void run() {
                    log.info("当前时间 {}", System.currentTimeMillis());
                }
            });
        }
    }

    @Test
    public void testFuture() {
        List<FutureTask<String>> futureTaskList = new ArrayList<>();

        FutureTask futureTask1 = new FutureTask<>(() -> {
            Thread.sleep(2000);// 模拟超时操作
            return "鸡翅";
        });
        FutureTask futureTask2 = new FutureTask<>(() -> {
            return "经典";
        });

        futureTaskList.add(futureTask1);
        futureTaskList.add(futureTask2);

        mailThreadPool.submit(futureTask1);
        mailThreadPool.submit(futureTask2);

        for (int i = 0; i < futureTaskList.size(); i++) {
            String name = CompletableFutureUtils.getResult(futureTaskList.get(i),
                    1, TimeUnit.SECONDS, "默认值", log);
            log.info("testFuture: {}", name);
        }
    }

    @Test
    public void publishEvent() {
        Person person = new Person();
        person.setName("鸡翅");
        person.setAge(18);
        personEventService.createPerson(person);
    }

    @Test
    public void pushTask() throws ParseException {
        SimpleDateFormat simpleDateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        MassMailTask massMailTask = new MassMailTask();
        massMailTask.setTaskId(1L);
        massMailTask.setStartTime(simpleDateFormat.parse("2025-01-31 23:07:00"));
        massMailTaskService.pushMassMailTaskQueue(massMailTask);
    }

    @Test
    public void deal() {
        String lockKey = "test.delay.task";
        String requestId = UUID.randomUUID().toString();
        try {
            boolean lock = redisShareLockUtil.lock(lockKey, requestId, 5L);
            if (!lock) {
                return;
            }
            Set<Long> taskIdSet = massMailTaskService.popMassMailTaskQueue();
            log.info("定时任务拉取: {}", JSON.toJSONString(taskIdSet));
            // todo 其他业务逻辑
        } catch (Exception e) {
            log.error("任务拉取异常: {}", e.getMessage(), e);
        } finally {
            redisShareLockUtil.unlock(lockKey, requestId);
        }
    }
}
