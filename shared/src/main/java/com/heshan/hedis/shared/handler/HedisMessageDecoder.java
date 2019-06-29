package com.heshan.hedis.shared.handler;

import com.heshan.hedis.shared.codec.*;
import com.heshan.hedis.shared.exception.HedisProtocolException;
import com.heshan.hedis.shared.util.HedisMessageUtils;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.ByteToMessageDecoder;
import io.netty.util.CharsetUtil;

import java.util.LinkedList;
import java.util.List;

/**
 * HedisMessageDecoder
 *
 * @author heshan
 * @date 2019-06-28
 */
public class HedisMessageDecoder extends ByteToMessageDecoder {

    private Reader reader;

    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) throws Exception {
        for (; ; ) {
            if (in.readableBytes() <= 0) {
                return;
            }

            if (reader == null) {
                reader = readReader(in);
            }

            HedisMessage msg = reader.read(in);
            if (msg != null) {
                out.add(msg);
                reader = null;
            } else {
                return;
            }
        }
    }

    private Reader readReader(ByteBuf buf) {
        char ch = (char) buf.readByte();
        switch (ch) {
            case ':':
                return new NumberMessageReader();
            case '*':
                return new ArrayMessageReader();
            case '+':
                return new StringMessageReader();
            case '-':
                return new ErrorMessageReader();
            case '$':
                return new BatchMessageReader();
            default:
                throw new HedisProtocolException();
        }
    }

    private interface Reader<T extends HedisMessage> {

        T read(ByteBuf buf);
    }

    private abstract class AbstractReader<T extends HedisMessage> implements Reader<T> {

        protected boolean finish;

        protected T message;

        @Override
        public T read(ByteBuf buf) {
            if (!finish) {
                doRead(buf);
            }

            return finish ? message : null;
        }

        protected abstract void doRead(ByteBuf buf);
    }

    private class StringMessageReader extends AbstractReader<StringHedisMessage> {

        @Override
        protected void doRead(ByteBuf buf) {
            ByteBuf line = HedisMessageUtils.readLine(buf);
            if (line == null) {
                return;
            }

            String content = line.readCharSequence(line.readableBytes(), CharsetUtil.UTF_8).toString();
            message = new StringHedisMessage(content);
            finish = true;
        }
    }

    private class ErrorMessageReader extends AbstractReader<ErrorHedisMessage> {

        @Override
        protected void doRead(ByteBuf buf) {
            ByteBuf line = HedisMessageUtils.readLine(buf);
            if (line == null) {
                return;
            }

            String content = line.readCharSequence(line.readableBytes(), CharsetUtil.UTF_8).toString();
            message = new ErrorHedisMessage(content);
            finish = true;
        }
    }

    private class NumberMessageReader extends AbstractReader<NumberHedisMessage> {

        @Override
        protected void doRead(ByteBuf buf) {
            ByteBuf line = HedisMessageUtils.readLine(buf);
            if (line == null) {
                return;
            }

            long value = HedisMessageUtils.readNumber(line);
            message = new NumberHedisMessage(value);
            finish = true;
        }
    }

    private class BatchMessageReader extends AbstractReader<BatchHedisMessage> {

        private int length = -2;

        @Override
        protected void doRead(ByteBuf buf) {
            if (length < -1) {
                ByteBuf line = HedisMessageUtils.readLine(buf);
                if (line == null) {
                    return;
                }

                length = (int) HedisMessageUtils.readNumber(line);
                if (length < -1) {
                    throw new HedisProtocolException();
                }
            }

            if (length == -1) {
                finish = true;
                return;
            }

            ByteBuf data = HedisMessageUtils.readLength(buf, length);
            if (data == null) {
                return;
            }

            String content = data.readCharSequence(data.readableBytes(), CharsetUtil.UTF_8).toString();
            message = new BatchHedisMessage(content);
            finish = true;
        }
    }

    private class ArrayMessageReader extends AbstractReader<ArrayHedisMessage> {

        private LinkedList<HedisMessage> messages = new LinkedList<>();

        private int length = -1;

        private Reader reader = null;

        @Override
        protected void doRead(ByteBuf buf) {
            readLength(buf);

            while (buf.readableBytes() > 0) {
                if (reader != null) {
                    HedisMessage msg = reader.read(buf);
                    if (msg == null) {
                        return;
                    }

                    messages.add(msg);
                    reader = null;

                    if (messages.size() == length) {
                        finish = true;
                        message = new ArrayHedisMessage(messages);
                        return;
                    }
                } else {
                    reader = readReader(buf);
                }
            }
        }

        private void readLength(ByteBuf buf) {
            if (length >= 0) {
                return;
            }

            ByteBuf line = HedisMessageUtils.readLine(buf);
            if (line == null) {
                return;
            }
            int value = (int) HedisMessageUtils.readNumber(line);
            if (value < 0) {
                throw new HedisProtocolException();
            }

            length = value;
        }
    }
}
