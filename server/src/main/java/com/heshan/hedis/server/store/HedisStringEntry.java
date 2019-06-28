package com.heshan.hedis.server.store;

/**
 * HedisStringEntry
 *
 * @author heshan
 * @date 2019-06-28
 */
public class HedisStringEntry extends HedisEntry {

    private String value;

    public HedisStringEntry(String value) {
        super();

        this.value = value;
    }

    public void set(String value) {
        this.value = value;
    }

    public String get() {
        return value;
    }
}
