package com.heshan.hedis.server.command.connection;

import com.heshan.hedis.server.command.AbstractHedisCommand;
import com.heshan.hedis.server.command.HedisCommandArgument;
import com.heshan.hedis.shared.codec.StringHedisMessage;

public class PingCommand extends AbstractHedisCommand {

    @Override
    public void execute(HedisCommandArgument arg) {
        StringHedisMessage res = new StringHedisMessage("PONG");
        arg.session().writeAndFlush(res);
    }
}
