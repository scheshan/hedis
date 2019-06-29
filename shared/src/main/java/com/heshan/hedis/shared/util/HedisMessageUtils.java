package com.heshan.hedis.shared.util;

import com.heshan.hedis.shared.codec.*;
import com.heshan.hedis.shared.exception.HedisProtocolException;
import io.netty.buffer.ByteBuf;
import io.netty.util.ByteProcessor;

/**
 * HedisMessageUtils
 *
 * @author heshan
 * @date 2019-06-28
 */
public final class HedisMessageUtils {

    private static final byte[] CRLF = new byte[]{'\r', '\n'};

    public static ByteBuf readLine(ByteBuf buf) {
        if (buf.readableBytes() <= 0) {
            return null;
        }

        int index = buf.forEachByte(ByteProcessor.FIND_LF);
        if (index < 0) {
            return null;
        } else if (index < 1) {
            throw new HedisProtocolException();
        }

        int length = index - buf.readerIndex() - 1;
        ByteBuf line = buf.slice(buf.readerIndex(), length);
        buf.readerIndex(buf.readerIndex() + length);
        readCRLF(buf);

        return line;
    }

    public static void readCRLF(ByteBuf buf) {
        if (buf.readableBytes() < 2) {
            throw new HedisProtocolException();
        }

        byte[] data = new byte[2];
        buf.readBytes(data);

        if (data[0] != CRLF[0] || data[1] != CRLF[1]) {
            throw new HedisProtocolException();
        }
    }

    public static long readNumber(ByteBuf buf) {
        long result = 0;
        boolean positive = true;
        char first = (char) buf.readByte();
        if (first == '-') {
            positive = false;
        } else {
            buf.readerIndex(buf.readerIndex() - 1);
        }

        for (int i = 0; i < buf.readableBytes(); i++) {
            char ch = (char) buf.readByte();
            if (ch < '0' && ch > '9') {
                throw new HedisProtocolException();
            }

            result = result * 10 + (ch - '0');
        }

        return positive ? result : -1 * result;
    }

    public static ByteBuf readLength(ByteBuf buf, int length) {
        if (buf.readableBytes() < length + 2) {
            return null;
        }

        ByteBuf content = buf.slice(buf.readerIndex(), length);
        buf.readerIndex(buf.readerIndex() + length);
        readCRLF(buf);

        return content;
    }

    public static AbstractHedisMessage readMessage(ByteBuf buf) {
        AbstractHedisMessage msg;

        char ch = (char) buf.readByte();
        switch (ch) {
            case ':':
                msg = new NumberHedisMessage();
                break;
            case '-':
                msg = new ErrorHedisMessage();
                break;
            case '+':
                msg = new StringHedisMessage();
                break;
            case '$':
                msg = new BatchHedisMessage();
                break;
            case '*':
                msg = new ArrayHedisMessage();
                break;
            default:
                throw new HedisProtocolException();
        }

        return msg;
    }

    public static void writeChar(ByteBuf buf, char ch) {
        buf.writeByte((byte) ch);
    }

    public static void writeCRLF(ByteBuf buf) {
        buf.writeBytes(CRLF);
    }
}
