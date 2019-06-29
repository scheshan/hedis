package com.heshan.hedis.server.command;

import com.heshan.hedis.shared.codec.BatchHedisMessage;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class SingleThreadCommandInvoker implements CommandInvoker {

    private ExecutorService executorService = Executors.newSingleThreadExecutor();

    @Override
    public void enqueue(RequestCommand cmd) {
        executorService.submit(new CommandTask(cmd));
    }

    private class CommandTask implements Runnable {

        private RequestCommand command;

        public CommandTask(RequestCommand command) {
            this.command = command;
        }

        @Override
        public void run() {
            command.session().writeAndFlush(new BatchHedisMessage(null));
        }
    }
}
