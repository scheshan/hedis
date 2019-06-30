package com.heshan.hedis.server.command;

import com.heshan.hedis.server.command.connection.PingCommand;

import java.util.HashMap;
import java.util.Map;

public class CommandManager {

    private Map<String, HedisCommand> commandMap = new HashMap<>();

    private CommandManager() {
        commandMap.put("ping", new PingCommand());
    }

    public HedisCommand createCommand(String name) {
        HedisCommand cmd = commandMap.get(name.toLowerCase());
        return cmd;
    }

    private static CommandManager instance = new CommandManager();

    public static CommandManager getInstance() {
        return instance;
    }
}
