package com.heshan.hedis.server.command;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class Executor {

    private static ExecutorService executorService = Executors.newSingleThreadExecutor();

    private static CommandFactory commandFactory = CommandFactory.getInstance();

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
            commandFactory.createCommand("").execute(request.session(), null);
        }
    }
}
