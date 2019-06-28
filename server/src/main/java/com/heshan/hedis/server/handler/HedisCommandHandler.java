package com.heshan.hedis.server.handler;

import com.heshan.hedis.shared.codec.AbstractHedisMessage;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.util.CharsetUtil;

/**
 * HedisCommandHandler
 *
 * @author heshan
 * @date 2019-06-28
 */
public class HedisCommandHandler extends ChannelInboundHandlerAdapter {

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        AbstractHedisMessage message = (AbstractHedisMessage) msg;
        System.out.println(message.content());

        ByteBuf buf = ctx.alloc().buffer();
        buf.writeCharSequence("*-1\r\n", CharsetUtil.UTF_8);
        ctx.writeAndFlush(buf);
    }
}
