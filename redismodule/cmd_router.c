#include "redismodule.h"
#include "report.h"

#include <stdlib.h>
#include <string.h>
#include <ffi.h>

static int router_command(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    RedisModule_AutoMemory(ctx);

    if (argc < 2) {
        RedisModule_WrongArity(ctx);
        return REDISMODULE_OK;
    }

    RedisModuleCallReply *reply;
    size_t len;
    const char *realCmd = RedisModule_StringPtrLen(argv[1], &len);

    if (argc == 2) {
        reply = RedisModule_Call(ctx, realCmd, "");
    } else {
        char *fmt = (char *) malloc(argc-1);
        memset(fmt, 's', argc-2);
        fmt[argc-2] = '\0';

        ffi_cif cif;
        ffi_type **ffi_argv = (ffi_type **)malloc(sizeof(ffi_type *) * (argc + 1));
        void *result;
        int i;
        for (i = 0; i < argc+1; i++) {
            ffi_argv[i] = &ffi_type_pointer;
        }
        void **values = (void **) malloc(sizeof(void *) * (argc + 1));
        values[0] = &ctx;
        values[1] = &realCmd;
        values[2] = &fmt;
        for (i = 3; i < argc + 1; i++) {
            values[i] = &argv[i-1];
        }

        if (ffi_prep_cif(&cif, FFI_DEFAULT_ABI, (argc + 1), &ffi_type_pointer, ffi_argv) == FFI_OK) {
            ffi_call(&cif, FFI_FN(RedisModule_Call), &result, values);
        }
        reply = (RedisModuleCallReply *) result;
    }
    if (reply == NULL) {
        const char *err = "command error";
        RedisModule_ReplyWithError(ctx, err);
    } else {
        report(ctx);
        RedisModule_ReplyWithCallReply(ctx, reply);
    }
    return REDISMODULE_OK;
}


int RedisModule_OnLoad(RedisModuleCtx *ctx) {
    if (RedisModule_Init(ctx, "rt", 1, REDISMODULE_APIVER_1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx, "rt.s", router_command,
                "no-monitor fast", 1, 1, 1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    return REDISMODULE_OK;
}
