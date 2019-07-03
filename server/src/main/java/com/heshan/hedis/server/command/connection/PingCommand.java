package com.heshan.hedis.server.command.connection;

import com.heshan.hedis.server.command.AbstractHedisCommand;
import com.heshan.hedis.server.command.HedisCommandArgument;
import com.heshan.hedis.shared.codec.StringHedisMessage;

public class PingCommand extends AbstractHedisCommand {

    private static final String PONG = "PONG";

    public PingCommand() {
        super(1, 0);
    }

    @Override
    public void doExecute(HedisCommandArgument arg) {
        String reply;
        if (arg.args().length > 0) {
            reply = arg.args()[0];
        } else {
            reply = PONG;
        }

        StringHedisMessage res = new StringHedisMessage(reply);
        arg.session().writeAndFlush(res);
    }
}
