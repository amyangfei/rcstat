#include "redismodule.h"
#include <string.h>
#include <nanomsg/nn.h>
#include <nanomsg/pipeline.h>

#define RCSTAT_CONN_KEY "rc_stat_conn"
#define RCSTAT_CONN_FD "fd"
#define RCSTAT_CONN_ADDR "addr"
#define RCSTAT_CONN_KEY_LEN 12
#define RCSTAT_CONN_FD_LEN 2
#define RCSTAT_CONN_ADDR_LEN 4


void error_handler() {
}

void report(RedisModuleCtx *ctx) {
    /* Open key and verify it is empty or a zset. */
    RedisModuleString *rms = RedisModule_CreateString(ctx, RCSTAT_CONN_KEY, RCSTAT_CONN_KEY_LEN);
    RedisModuleKey *key = RedisModule_OpenKey(ctx, rms, REDISMODULE_READ|REDISMODULE_WRITE);
    int key_type = RedisModule_KeyType(key);
    if (key_type == REDISMODULE_KEYTYPE_EMPTY || key_type != REDISMODULE_KEYTYPE_HASH) {
        return error_handler();
    }
    RedisModuleString *fd;
    RedisModuleString *fd_key = RedisModule_CreateString(ctx, RCSTAT_CONN_FD, RCSTAT_CONN_FD_LEN);
    RedisModule_HashGet(key, REDISMODULE_HASH_NONE, fd_key, &fd, NULL);
    if (!fd) {
        RedisModuleString *addr;
        RedisModuleString *addr_key = RedisModule_CreateString(ctx, RCSTAT_CONN_ADDR, RCSTAT_CONN_ADDR_LEN);
        RedisModule_HashGet(key, REDISMODULE_HASH_NONE, addr_key, &addr, NULL);
        if (addr) {
            size_t len;
            const char *addr_str = RedisModule_StringPtrLen(addr, &len);
            // create new connection
            int sock = nn_socket(AF_SP, NN_PUSH);
            if (sock < 0) {
                return error_handler();
            }
            int ret = nn_connect(sock, addr_str);
            if (ret < 0) {
                return error_handler();
            }
            fd = RedisModule_CreateStringFromLongLong(ctx, sock);
            RedisModule_HashSet(key, REDISMODULE_HASH_NONE, fd_key, fd, NULL);
            /* nn_shutdown(sock, 0); */
        }
    }
    long long fd_int;
    RedisModule_StringToLongLong(fd, &fd_int);
    const char *msg = "hello from reporter";
    nn_send(fd_int, msg, strlen(msg), 0);
}
