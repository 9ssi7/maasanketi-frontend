package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type queryTracer struct {
	log *zap.SugaredLogger
}

func (tracer *queryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.log.Debugw("Executing command", "sql", data.SQL, "args", data.Args)

	return ctx
}

func (tracer *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {

}
