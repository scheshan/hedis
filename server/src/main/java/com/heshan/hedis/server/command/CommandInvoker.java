package com.heshan.hedis.server.command;

public interface CommandInvoker {

    void enqueue(RequestCommand cmd);

    static CommandInvoker instance() {
        return new SingleThreadCommandInvoker();
    }
}
