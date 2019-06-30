package com.heshan.hedis.server.command;

import com.heshan.hedis.server.command.connection.PingCommand;

import java.util.HashMap;
import java.util.Map;

public class CommandFactory {

    private Map<String, HedisCommand> commandMap = new HashMap<>();

    private CommandFactory() {
        commandMap.put("ping", new PingCommand());
    }

    public HedisCommand createCommand(String name) {
        HedisCommand cmd = commandMap.get(name.toLowerCase());
        return cmd;
    }

    private static CommandFactory instance = new CommandFactory();

    public static CommandFactory getInstance() {
        return instance;
    }
}
