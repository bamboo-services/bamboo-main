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

/**
 * LinkGetAdminEntity
 *
 * 用于获取所有的链接，该接口只有管理员可以进行操作，其他用户不应该通过该接口进行数据的获取操作
 *
 * @param links 友链列表
 * @param total 友链总数
 */
export type LinkGetAdminEntity = {
    links: LinkDO[];
    total: number;
}

/**
 * LinkDO
 *
 * 友链数据对象
 *
 * @param id 友链ID
 * @param webmaster_email 网站管理员邮箱
 * @param service_provider 服务提供商
 * @param site_name 网站名称
 * @param site_url 网站URL
 * @param site_logo 网站Logo
 * @param cdn_logo_url CDN Logo URL
 * @param site_description 网站描述
 * @param site_rss_url 网站RSS URL
 * @param has_adv 是否有广告
 * @param desired_location 期望位置
 * @param location 位置
 * @param desired_color 期望颜色
 * @param color 颜色
 * @param webmaster_remark 网站管理员备注
 * @param remark 备注
 * @param status 状态
 * @param able_connect 是否可以连接
 * @param created_at 创建时间
 * @param updated_at 更新时间
 * @param deleted_at 删除时间
 */
export type LinkDO = {
    id: number;
    webmaster_email: string;
    service_provider: string;
    site_name: string;
    site_url: string;
    site_logo: string;
    cdn_logo_url: string;
    site_description: string;
    site_rss_url: string;
    has_adv: boolean;
    desired_location: number;
    location: number;
    desired_color: number;
    color: number;
    webmaster_remark: string;
    remark: string;
    status: number;
    able_connect: boolean;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
}
