package com.heshan.hedis.server.handler;

import com.heshan.hedis.server.command.Executor;
import com.heshan.hedis.server.command.RequestWrapper;
import com.heshan.hedis.server.session.SessionManager;
import com.heshan.hedis.shared.codec.HedisMessage;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;

/**
 * CommandHandler
 *
 * @author heshan
 * @date 2019-06-28
 */
public class CommandHandler extends ChannelInboundHandlerAdapter {

    private static Executor executor = Executor.getInstance();

    private static SessionManager sessionManager = SessionManager.getInstance();

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        HedisMessage message = (HedisMessage) msg;

        RequestWrapper request = new RequestWrapper(message, sessionManager.get(ctx.channel()));
        executor.execute(request);
    }
}
