package com.heshan.hedis.server.command.ping;

import com.heshan.hedis.server.command.AbstractHedisCommand;
import com.heshan.hedis.server.session.HedisSession;
import com.heshan.hedis.shared.codec.StringHedisMessage;

public class PingCommand extends AbstractHedisCommand {

    @Override
    public void execute(HedisSession session, String[] args) {
        StringHedisMessage res = new StringHedisMessage("PONG");
        session.writeAndFlush(res);
    }
}
