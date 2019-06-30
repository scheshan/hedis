package com.heshan.hedis.server.command;

import com.heshan.hedis.server.session.HedisSession;
import com.heshan.hedis.shared.codec.ArrayHedisMessage;
import com.heshan.hedis.shared.codec.ErrorHedisMessage;
import com.heshan.hedis.shared.codec.HedisMessage;
import com.heshan.hedis.shared.exception.HedisProtocolException;

import java.util.Iterator;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class Executor {

    private static ExecutorService executorService = Executors.newSingleThreadExecutor();

    private static CommandFactory commandFactory = CommandFactory.getInstance();

    private Executor() {

    }

    public void execute(HedisSession session, HedisMessage message) {
        executorService.submit(new ExecutorTask(session, message));
    }

    private static Executor instance = new Executor();

    public static Executor getInstance() {
        return instance;
    }

    private class ExecutorTask implements Runnable {

        private HedisSession session;

        private HedisMessage message;

        private String commandName;

        private String[] commandArgs;

        public ExecutorTask(HedisSession session, HedisMessage message) {
            this.session = session;
            this.message = message;
        }

        @Override
        public void run() {
            try {
                parse();
                HedisCommand command = commandFactory.createCommand(commandName);
                if (command == null) {
                    throw new HedisProtocolException();
                }
                command.execute(session, commandArgs);
            } catch (Exception ex) {
                ErrorHedisMessage res = new ErrorHedisMessage(ex.getMessage());
                session.writeAndFlush(res);
            }
        }

        private void parse() {
            if (message instanceof ArrayHedisMessage) {
                Iterator<HedisMessage> messages = ((ArrayHedisMessage) message).messages().iterator();
                commandName = messages.next().content();

                commandArgs = new String[((ArrayHedisMessage) message).size() - 1];
                int i = 0;
                while (messages.hasNext()) {
                    commandArgs[i++] = messages.next().content();
                }
            } else {
                commandName = message.content();
                commandArgs = new String[0];
            }
        }
    }
}
