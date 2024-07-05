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
 * LocationGetAdminEntity
 *
 * 用于获取所有的位置，该接口只有管理员可以进行操作，其他用户不应该通过该接口进行数据的获取操作
 *
 * @param locations 位置列表
 * @param total 位置总数
 */
export type LocationGetAdminEntity = {
    locations: LocationDO[];
    total: number;
}

/**
 * LocationDO
 *
 * 位置数据对象
 *
 * @param id 位置ID
 * @param sort 位置排序
 * @param name 位置名称
 * @param display_name 位置显示名称
 * @param description 位置描述
 * @param reveal 是否显示
 * @param created_at 创建时间
 * @param updated_at 更新时间
 */
export type LocationDO = {
    id: number;
    sort: number;
    name: string;
    display_name: string;
    description: string;
    reveal: boolean;
    created_at: string;
    updated_at: string;
}
