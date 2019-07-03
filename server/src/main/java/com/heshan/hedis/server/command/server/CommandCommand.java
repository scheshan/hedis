package com.heshan.hedis.server.command.server;

import com.heshan.hedis.server.command.AbstractHedisCommand;
import com.heshan.hedis.server.command.HedisCommandArgument;
import com.heshan.hedis.shared.codec.StringHedisMessage;

/**
 * CommandCommand
 *
 * @author heshan
 * @date 2019-07-03
 */
public class CommandCommand extends AbstractHedisCommand {

    public CommandCommand() {
        super(0, 0);
    }

    @Override
    protected void doExecute(HedisCommandArgument arg) {
        StringHedisMessage res = new StringHedisMessage("OK");
        arg.session().writeAndFlush(res);
    }
}
