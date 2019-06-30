package com.heshan.hedis.server.command;

public interface HedisCommand {

    void execute(HedisCommandArgument arg);
}
