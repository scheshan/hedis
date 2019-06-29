package com.heshan.hedis.shared.handler;

import com.heshan.hedis.shared.codec.*;
import com.heshan.hedis.shared.util.HedisMessageUtils;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToByteEncoder;
import io.netty.util.CharsetUtil;

/**
 * HedisMessageEncoder
 *
 * @author heshan
 * @date 2019-06-29
 */
public class HedisMessageEncoder extends MessageToByteEncoder<HedisMessage> {

    @Override
    protected void encode(ChannelHandlerContext ctx, HedisMessage msg, ByteBuf out) throws Exception {
        writeMessage(out, msg);
    }

    private void writeMessage(ByteBuf buf, HedisMessage msg) {
        if (msg instanceof StringHedisMessage) {
            writeStringMessage(buf, (StringHedisMessage) msg);
        } else if (msg instanceof NumberHedisMessage) {
            writeNumberMessage(buf, (NumberHedisMessage) msg);
        } else if (msg instanceof ErrorHedisMessage) {
            writeErrorMessage(buf, (ErrorHedisMessage) msg);
        } else if (msg instanceof BatchHedisMessage) {
            writeBatchMessage(buf, (BatchHedisMessage) msg);
        } else if (msg instanceof ArrayHedisMessage) {
            writeArrayMessage(buf, (ArrayHedisMessage) msg);
        } else {
            throw new UnsupportedOperationException();
        }
    }

    private void writeNumberMessage(ByteBuf buf, NumberHedisMessage msg) {
        HedisMessageUtils.writeChar(buf, ':');
        writeNumber(buf, msg.value());
    }

    private void writeStringMessage(ByteBuf buf, StringHedisMessage msg) {
        HedisMessageUtils.writeChar(buf, '+');
        writeString(buf, msg.content());
    }

    private void writeErrorMessage(ByteBuf buf, ErrorHedisMessage msg) {
        HedisMessageUtils.writeChar(buf, '-');
        writeString(buf, msg.content());
    }

    private void writeBatchMessage(ByteBuf buf, BatchHedisMessage msg) {
        HedisMessageUtils.writeChar(buf, '$');
        String content = msg.content();
        if (content == null) {
            writeNumber(buf, -1);
        } else {
            writeNumber(buf, content.length());
            writeString(buf, content);
        }
    }

    private void writeArrayMessage(ByteBuf buf, ArrayHedisMessage msg) {
        HedisMessageUtils.writeChar(buf, '*');
        writeNumber(buf, msg.size());
        for (HedisMessage child : msg.messages()) {
            writeMessage(buf, child);
        }
    }

    private void writeNumber(ByteBuf buf, long number) {
        buf.writeCharSequence(String.valueOf(number), CharsetUtil.UTF_8);
        HedisMessageUtils.writeCRLF(buf);
    }

    private void writeString(ByteBuf buf, String content) {
        buf.writeCharSequence(content, CharsetUtil.UTF_8);
        HedisMessageUtils.writeCRLF(buf);
    }
}
