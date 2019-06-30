package com.heshan.hedis.server.command;

import com.heshan.hedis.server.session.SessionManager;

public abstract class AbstractHedisCommand implements HedisCommand {

    protected final CommandManager commandManager = CommandManager.getInstance();

    protected final SessionManager sessionManager = SessionManager.getInstance();
}
