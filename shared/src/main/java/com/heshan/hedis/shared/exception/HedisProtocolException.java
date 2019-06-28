package com.heshan.hedis.shared.exception;

/**
 * HedisProtocolException
 *
 * @author heshan
 * @date 2019-06-28
 */
public class HedisProtocolException extends RuntimeException {

    public HedisProtocolException() {
        super();
    }

    public HedisProtocolException(String message) {
        super(message);
    }
}
