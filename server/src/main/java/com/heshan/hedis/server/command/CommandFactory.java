package com.heshan.hedis.server.command;

import com.heshan.hedis.server.command.ping.PingCommand;

public class CommandFactory {

    private CommandFactory() {

    }

    public void init() {

    }

    public HedisCommand createCommand(String name) {
        return new PingCommand();
    }

    private static CommandFactory instance = new CommandFactory();

    public static CommandFactory getInstance() {
        return instance;
    }
}
