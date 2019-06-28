package com.heshan.hedis.server.store;

import java.util.HashMap;
import java.util.Map;

/**
 * HedisStore
 *
 * @author heshan
 * @date 2019-06-28
 */
public class HedisStore {

    private Map<String, HedisEntry> entryMap = new HashMap<>();

    public HedisEntry get(String key) {
        return entryMap.get(key);
    }

    public void set(String key, HedisEntry entry) {
        entryMap.put(key, entry);
    }

    public boolean exist(String key) {
        return entryMap.containsKey(key);
    }

    public void remove(String key) {
        entryMap.remove(key);
    }

    public int size() {
        return entryMap.size();
    }
}
