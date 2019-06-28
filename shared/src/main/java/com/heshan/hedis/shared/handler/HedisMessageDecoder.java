package com.heshan.hedis.shared.handler;

import com.heshan.hedis.shared.codec.AbstractHedisMessage;
import com.heshan.hedis.shared.util.HedisMessageUtils;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.ByteToMessageDecoder;

import java.util.List;

/**
 * HedisMessageDecoder
 *
 * @author heshan
 * @date 2019-06-28
 */
public class HedisMessageDecoder extends ByteToMessageDecoder {

    private AbstractHedisMessage msg;

    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) throws Exception {
        for (; ; ) {
            if (in.readableBytes() <= 0) {
                return;
            }

            if (msg == null) {
                msg = HedisMessageUtils.readMessage(in);
            }

            msg.read(in);

            if (msg.isFinish()) {
                out.add(msg);
            } else {
                return;
            }
        }
    }
}
