package boot

import (
	"context"
	"develop/internal/dao"
	"develop/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/glog"
)

func InitialDatabase(ctx context.Context) {
	var err error

	/*
	 * 检查数据表是否完善
	 */
	if _, err := dao.XfIndex.DB().Exec(ctx, "SELECT * FROM xf_index"); err != nil {
		// 创建数据表
		err := dao.XfIndex.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 创建 xf_index 表
			if _, err := tx.Exec("CREATE TABLE xf_index ( id serial CONSTRAINT xf_index_pk PRIMARY KEY, key varchar(40) NOT NULL, value text, created_at timestamp NOT NULL DEFAULT now(), updated_at timestamp );"); err != nil {
				return err
			}
			// 创建 COMMIT
			if _, err := tx.Exec("comment on table xf_index is '筱主页信息数据表';comment on column xf_index.id is '主键';comment on column xf_index.key is '键';comment on column xf_index.value is '值';comment on column xf_index.created_at is '创建时间';comment on column xf_index.updated_at is '修改时间';"); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic("[BOOT] 数据表 xf_index 创建失败")
		} else {
			glog.Debug(ctx, "[BOOT] 数据表 xf_index 创建成功")
		}
	} else {
		glog.Debug(ctx, "[BOOT] 数据表 xf_index 已存在")
	}

	/**
	 * 检查数据表信息是否完整
	 */
	if _, err = dao.XfIndex.DB().Exec(ctx, "SELECT * FROM xf_index WHERE `key` = 'version'"); err != nil {
		// 插入版本信息
		if _, err = dao.XfIndex.Ctx(ctx).Data(entity.XfIndex{Key: "version", Value: "v1.0.0"}).Insert(); err != nil {
			glog.Error(ctx, "[BOOT] 数据表 xf_index 中 version 数据创建失败")
		} else {
			glog.Debug(ctx, "[BOOT] 数据表 xf_index 中 version 数据创建成功")
		}
		//
	} else {
		glog.Debug(ctx, "[BOOT] 数据表 xf_index 中 version 数据已存在")
	}
}
