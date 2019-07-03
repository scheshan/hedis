package com.heshan.hedis.server.command.connection;

import com.heshan.hedis.server.command.AbstractHedisCommand;
import com.heshan.hedis.server.command.HedisCommandArgument;
import com.heshan.hedis.shared.codec.StringHedisMessage;

/**
 * EchoCommand
 *
 * @author heshan
 * @date 2019-06-30
 */
public class EchoCommand extends AbstractHedisCommand {

    public EchoCommand() {
        super(1, 1);
    }

    @Override
    public void doExecute(HedisCommandArgument arg) {
        StringHedisMessage res = new StringHedisMessage(arg.args()[0]);
        arg.session().writeAndFlush(res);
    }
}
