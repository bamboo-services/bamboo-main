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

create table xf_link_list
(
    id               bigserial
        constraint xf_link_list_pk
            primary key,
    webmaster_email  varchar,
    service_provider varchar                         not null,
    site_name        varchar                         not null,
    site_url         varchar                         not null,
    site_logo        varchar                         not null,
    site_description varchar                         not null,
    site_rss_url     varchar,
    has_adv          boolean   default false         not null,
    desired_location integer   default 0             not null,
    location         integer   default '-1'::integer not null,
    desired_color    integer   default 0             not null,
    color            integer   default '-1'::integer not null,
    webmaster_remark varchar,
    remark           varchar,
    status           smallint  default 0             not null,
    created_at       timestamp default now()         not null,
    updated_at       timestamp,
    deleted_at       timestamp
);

comment on table xf_link_list is '链接列表';
comment on column xf_link_list.id is '主键';
comment on column xf_link_list.webmaster_email is '站长邮箱';
comment on column xf_link_list.service_provider is '服务提供商';
comment on column xf_link_list.site_name is '站点名字';
comment on column xf_link_list.site_url is '站点地址';
comment on column xf_link_list.site_logo is '站点 logo';
comment on column xf_link_list.site_description is '站点描述';
comment on column xf_link_list.site_rss_url is '站点订阅地址';
comment on column xf_link_list.has_adv is '是否有广告';
comment on column xf_link_list.desired_location is '理想位置';
comment on column xf_link_list.location is '所在位置';
comment on column xf_link_list.desired_color is '理想颜色';
comment on column xf_link_list.color is '颜色';
comment on column xf_link_list.webmaster_remark is '站长留言';
comment on column xf_link_list.remark is '我的留言';
comment on column xf_link_list.status is '0 待审核，1 通过，-1 审核拒绝';
comment on column xf_link_list.created_at is '创建时间';
comment on column xf_link_list.updated_at is '修改时间';
comment on column xf_link_list.deleted_at is '删除时间';