package com.heshan.hedis.server.command;

import com.heshan.hedis.server.command.connection.EchoCommand;
import com.heshan.hedis.server.command.connection.PingCommand;
import com.heshan.hedis.server.command.server.CommandCommand;

import java.util.HashMap;
import java.util.Map;

public class CommandManager {

    private Map<String, HedisCommand> commandMap = new HashMap<>();

    private CommandManager() {
        //region connection
        commandMap.put("ping", new PingCommand());
        commandMap.put("echo", new EchoCommand());
        //endregion

        //region server
        commandMap.put("command", new CommandCommand());
        //endregion
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
