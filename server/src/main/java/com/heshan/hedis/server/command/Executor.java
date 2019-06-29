package com.heshan.hedis.server.command;

import com.heshan.hedis.shared.codec.BatchHedisMessage;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class Executor {

    private ExecutorService executorService = Executors.newSingleThreadExecutor();

    private Executor() {

    }

    public void execute(RequestWrapper request) {
        executorService.submit(new ExecutorTask(request));
    }

    private static Executor instance = new Executor();

    public static Executor getInstance() {
        return instance;
    }

    private class ExecutorTask implements Runnable {

        private RequestWrapper request;

        public ExecutorTask(RequestWrapper request) {
            this.request = request;
        }

        @Override
        public void run() {
            request.session().writeAndFlush(new BatchHedisMessage(null));
        }
    }
}
