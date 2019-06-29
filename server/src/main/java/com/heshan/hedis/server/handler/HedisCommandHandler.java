package com.heshan.hedis.server.handler;

import com.heshan.hedis.server.command.CommandInvoker;
import com.heshan.hedis.server.command.RequestCommand;
import com.heshan.hedis.server.session.SessionManager;
import com.heshan.hedis.shared.codec.HedisMessage;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;

/**
 * HedisCommandHandler
 *
 * @author heshan
 * @date 2019-06-28
 */
public class HedisCommandHandler extends ChannelInboundHandlerAdapter {

    private static CommandInvoker invoker = CommandInvoker.instance();

    private static SessionManager sessionManager = SessionManager.getInstance();

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        HedisMessage message = (HedisMessage) msg;

        RequestCommand cmd = new RequestCommand(sessionManager.get(ctx.channel()), message);
        invoker.enqueue(cmd);
    }
}
