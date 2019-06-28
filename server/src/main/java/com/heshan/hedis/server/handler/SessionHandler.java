package com.heshan.hedis.server.handler;

import com.heshan.hedis.server.session.SessionManager;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;

/**
 * SessionHandler
 *
 * @author heshan
 * @date 2019-06-28
 */
public class SessionHandler extends ChannelInboundHandlerAdapter {

    private final SessionManager sessionManager = SessionManager.getInstance();

    @Override
    public void channelRegistered(ChannelHandlerContext ctx) throws Exception {
        sessionManager.sessionInit(ctx.channel());

        super.channelRegistered(ctx);
    }

    @Override
    public void channelUnregistered(ChannelHandlerContext ctx) throws Exception {
        sessionManager.sessionClose(ctx.channel());

        super.channelUnregistered(ctx);
    }
}
