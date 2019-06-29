package com.heshan.hedis.shared.codec;

import com.heshan.hedis.shared.exception.HedisProtocolException;
import com.heshan.hedis.shared.util.HedisMessageUtils;
import io.netty.buffer.ByteBuf;

import java.util.Collection;
import java.util.LinkedList;

/**
 * ArrayHedisMessage
 *
 * @author heshan
 * @date 2019-06-28
 */
public class ArrayHedisMessage extends AbstractHedisMessage {

    private LinkedList<AbstractHedisMessage> messages = new LinkedList<>();

    private int length = -1;

    @Override
    protected void doRead(ByteBuf buf) {
        readLength(buf);

        while (buf.readableBytes() > 0) {
            AbstractHedisMessage msg;
            if (!messages.isEmpty() && !messages.getLast().isFinish()) {
                msg = messages.getLast();
            } else {
                msg = HedisMessageUtils.readMessage(buf);
                messages.addLast(msg);
            }

            msg.read(buf);
            if (!msg.isFinish()) {
                return;
            }
            if (messages.size() == length) {
                finish = true;
                return;
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

    @Override
    public String content() {
        StringBuilder sb = new StringBuilder();
        for (AbstractHedisMessage msg : messages) {
            sb.append(msg.content);
            sb.append("\r\n");
        }

        return sb.toString();
    }

    public int length() {
        return length;
    }

    public Collection<AbstractHedisMessage> messages() {
        return messages;
    }
}
