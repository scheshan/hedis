package com.heshan.hedis.server.command;

import com.heshan.hedis.server.session.SessionManager;

public abstract class AbstractHedisCommand implements HedisCommand {

    protected final CommandManager commandManager = CommandManager.getInstance();

    protected final SessionManager sessionManager = SessionManager.getInstance();

    private int maxArgs;

    private int minArgs;

    protected AbstractHedisCommand(int maxArgs, int minArgs) {
        this.maxArgs = maxArgs;
        this.minArgs = minArgs;
    }

    protected abstract void doExecute(HedisCommandArgument arg);

    @Override
    public void execute(HedisCommandArgument arg) {
        int argLength = arg.args().length;

        if (minArgs > -1 && argLength < minArgs) {
            arg.session().writeError("Wrong argument number");
        }
        if (maxArgs > -1 && argLength > maxArgs) {
            arg.session().writeError("Wrong argument number");
        }

        doExecute(arg);
    }
}
