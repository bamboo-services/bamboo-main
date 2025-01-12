/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

-- auto-generated definition
create table xf_sponsor
(
    sponsor_uuid uuid                    not null
        constraint xf_sponsor_uuid_pk
            primary key,
    name         varchar(255)            not null,
    url          varchar,
    type         integer                 not null
        constraint xf_sponsor_xf_sponsor_type_id_fk
            references xf_sponsor_type
            on update cascade on delete cascade,
    money        double precision        not null,
    time         timestamp default now() not null,
    created_at   timestamp default now() not null,
    updated_at   timestamp default now() not null
);

comment on table xf_sponsor is '赞助表';
comment on column xf_sponsor.sponsor_uuid is '赞助主键';
comment on column xf_sponsor.name is '赞助者名称';
comment on column xf_sponsor.url is '地址';
comment on column xf_sponsor.type is '赞助方式';
comment on column xf_sponsor.money is '赞助金额';
comment on column xf_sponsor.time is '时间戳';
comment on column xf_sponsor.created_at is '创建时间';
comment on column xf_sponsor.updated_at is '修改时间';

