package com.jingdianjichi.tool;

import org.slf4j.Logger;

import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;

public class CompletableFutureUtils {
    public static <T> T getResult(Future<T> future, long timeOut, TimeUnit timeUnit, T defaultValue, Logger logger) {
        try {
            return future.get(timeOut, timeUnit);
        } catch (Exception e) {
            logger.error("CompletableFutureUtils.getResult.error:{}", e.getMessage(), e);
            return defaultValue;
        }
    }
}
