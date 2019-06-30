package com.heshan.hedis.server.handler;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;

public class DebugHandler extends ChannelInboundHandlerAdapter {

    private boolean debug;

    public DebugHandler(boolean debug) {
        this.debug = debug;
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        if (debug) {
            ByteBuf buf = (ByteBuf) msg;
            System.out.println(ByteBufUtil.prettyHexDump(buf));
        }

        super.channelRead(ctx, msg);
    }
}
