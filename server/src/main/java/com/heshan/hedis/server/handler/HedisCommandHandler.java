package com.heshan.hedis.server.handler;

import com.heshan.hedis.shared.codec.BatchHedisMessage;
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

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        HedisMessage message = (HedisMessage) msg;
        System.out.println(message.toString());

        HedisMessage res = new BatchHedisMessage(null);
        ctx.writeAndFlush(res);
    }
}
