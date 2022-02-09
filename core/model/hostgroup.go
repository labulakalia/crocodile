package model

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateHostgroup create hostgroup
func CreateHostgroup(ctx context.Context, name, remark, createByID string, hostids []string) error {
	createsql := `INSERT INTO crocodile_hostgroup (id,name,remark,createByID,hostIDs,createTime,updateTime) VALUES(?,?,?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, createsql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	createTime := time.Now().Unix()
	_, err = stmt.ExecContext(ctx,
		utils.GetID(),
		name,
		remark,
		createByID,
		strings.Join(hostids, ","),
		createTime,
		createTime)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// ChangeHostGroup change hostgroup
func ChangeHostGroup(ctx context.Context, hostids []string, id, remark string) error {
	changesql := `UPDATE crocodile_hostgroup SET hostIDs=?,remark=?,updateTime=? WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, changesql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx,
		strings.Join(hostids, ","),
		remark,
		time.Now().Unix(),
		id,
	)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// DeleteHostGroup delete hostgroup
func DeleteHostGroup(ctx context.Context, id string) error {
	sqldelete := `DELETE FROM crocodile_hostgroup WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, sqldelete)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// getHostGroups return hostgroup by id or hostgroupname
func getHostGroups(ctx context.Context, id, hgname string, limit, offset int) ([]define.HostGroup, int, error) {
	hgs := []define.HostGroup{}
	getsql := `SELECT 
					hg.id,
					hg.name,
					hg.remark,
					hg.hostIDs,
					hg.createByID,
					u.name,
					hg.createTime,
					hg.updateTime 
				FROM 
					crocodile_hostgroup as hg,crocodile_user as u
				WHERE
					hg.createByID = u.id`
	var count int
	args := []interface{}{}
	if id != "" {
		getsql += " AND hg.id=?"
		args = append(args, id)
	}
	if hgname != "" {
		getsql += " AND hg.name=?"
		args = append(args, hgname)
	}
	if limit > 0 {
		var err error
		count, err = countColums(ctx, getsql, args...)
		if err != nil {
			return hgs, 0, fmt.Errorf("countColums failed: %w", err)
		}
		getsql += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	conn, err := db.GetConn(ctx)
	if err != nil {
		return hgs, 0, fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return hgs, 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return hgs, 0, fmt.Errorf("stmt.QueryContext failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			hg                     define.HostGroup
			addrs                  string
			createTime, updateTime int64
		)
		err := rows.Scan(&hg.ID, &hg.Name, &hg.Remark,
			&addrs, &hg.CreateByUID, &hg.CreateBy, &createTime, &updateTime)
		if err != nil {
			log.Info("Scan result failed", zap.Error(err))
			continue
		}
		hg.HostsID = []string{}
		if addrs != "" {
			hg.HostsID = append(hg.HostsID, strings.Split(addrs, ",")...)

		}
		hg.CreateTime = utils.UnixToStr(createTime)
		hg.UpdateTime = utils.UnixToStr(updateTime)

		hgs = append(hgs, hg)
	}
	return hgs, count, nil
}

// GetHostGroups return all hostgroup
func GetHostGroups(ctx context.Context, limit, offset int) ([]define.HostGroup, int, error) {
	return getHostGroups(ctx, "", "", limit, offset)
}

// GetHostGroupByID return hostgroup by id
func GetHostGroupByID(ctx context.Context, id string) (*define.HostGroup, error) {
	hostgroups, _, err := getHostGroups(ctx, id, "", 0, 0)
	if err != nil {
		return nil, err
	}
	if len(hostgroups) != 1 {
		err = define.ErrNotExist{Value: id, Type: "hostgroup id"}
		return nil, err
	}
	return &hostgroups[0], nil
}

// GetHostsByHGID return hostgroup's host details
func GetHostsByHGID(ctx context.Context, hgid string) ([]*define.Host, error) {
	hostgroup, err := GetHostGroupByID(ctx, hgid)
	if err != nil {
		return nil, err
	}
	if len(hostgroup.HostsID) == 0 {
		return []*define.Host{}, nil
	}
	hosts, err := GetHostByIDS(ctx, hostgroup.HostsID)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

// GetHostGroupByName return hostgroup by name
func GetHostGroupByName(ctx context.Context, hgname string) (*define.HostGroup, error) {
	hostgroups, _, err := getHostGroups(ctx, "", hgname, 0, 0)
	if err != nil {
		return nil, err
	}

	if len(hostgroups) != 1 {
		return nil, errors.New("can not find hostgroup name: " + hgname)
	}
	return &hostgroups[0], nil
}

// RandHostID return execute worker ip
func RandHostID(hg *define.HostGroup) (string, error) {
	if len(hg.HostsID) == 0 {
		return "", errors.New("Can not find worker host")
	}
	hostid := hg.HostsID[rand.Int()%len(hg.HostsID)]
	return hostid, nil
}

// TODO V2

// CreateHostgroupv2 create new hostgroup
func CreateHostgroupv2(ctx context.Context, name, remark, createID string, hostids []string) error {

	hostgroup := &HostGroup{
		Name:     name,
		CreateID: createID,
		Hosts:    IDs(hostids),
		Remark:   remark,
	}

	err := gormdb.WithContext(ctx).Create(hostgroup).Error
	if err != nil {
		return fmt.Errorf("create hosrgroup %v failed: %w", hostgroup, err)
	}
	return nil
}

// ChangeHostGroupv2 change hostgroup
func ChangeHostGroupv2(ctx context.Context, hostids []string, id, remark string) error {
	hostgroup := &HostGroup{
		Hosts:  IDs(hostids),
		Remark: remark,
	}
	v, ok := ctx.Value("uid").(string)
	if ok {
		hostgroup.CurrentUID = v
	}

	res := gormdb.WithContext(ctx).Model(&HostGroup{}).Where("id = ?", id).Updates(hostgroup)
	if res.Error != nil {
		return fmt.Errorf("update hostgroup %v failed: %w", id, res.Error)
	}
	if res.RowsAffected == 0 {
		return define.ErrNotExist{Type: "hostgroup id", Value: id}
	}
	return nil
}

// DeleteHostGroupv2 delete hostgroup
func DeleteHostGroupv2(ctx context.Context, id string) error {
	hostgroup := HostGroup{
		Model: Model{
			ID: id,
		},
	}
	v, ok := ctx.Value("uid").(string)
	if ok {
		hostgroup.CurrentUID = v
	}

	res := gormdb.WithContext(ctx).Model(&HostGroup{}).Delete(&hostgroup)
	if res.Error != nil {
		return fmt.Errorf("delete hostgroup %s failed: %w", id, res.Error)
	}
	if res.RowsAffected == 0 {
		return define.ErrNotExist{Type: "hostgroup id", Value: id}
	}
	return nil
}

// GetHostGroupsv2 get all hostgroup
func GetHostGroupsv2(ctx context.Context, limit, offset int) ([]*HostGroup, int64, error) {
	var hgs = []*HostGroup{}
	var count int64
	err := gormdb.WithContext(ctx).Find(&hgs).Count(&count).Limit(limit).Offset(offset).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find hostgroup failed: %w", err)
	}

	return hgs, count, nil
}

// GetHostGroupByIDv2 get hostgroup id by id
func GetHostGroupByIDv2(ctx context.Context, id string) (*HostGroup, error) {
	var hg = HostGroup{}
	err := gormdb.WithContext(ctx).Where("id = ?", id).First(&hg).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("get hostgroup %s failed: %w", id, err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, define.ErrNotExist{Value: id, Type: "hostgroup id"}
	}
	return &hg, nil
}

// GetHostGroupByNamev2 get hostgroup by name
func GetHostGroupByNamev2(ctx context.Context, hgname string) (*HostGroup, error) {
	var hg = HostGroup{}
	err := gormdb.WithContext(ctx).Where("name = ?", hgname).First(&hg).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("get hostgroup %s failed: %w", hgname, err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, define.ErrNotExist{Value: hgname, Type: "hostgroup name"}
	}
	return &hg, nil
}

// GetHostsByHGIDv2 get hosts by hg id
func GetHostsByHGIDv2(ctx context.Context, hgid string) ([]*Host, error) {
	hg, err := GetHostGroupByIDv2(ctx, hgid)
	if err != nil {
		return nil, fmt.Errorf("get hostgroupd %s failed: %w", hgid, err)
	}
	hosts, err := GetHostsByIDSv2(ctx, hg.Hosts)
	if err != nil {
		return nil, fmt.Errorf("get hosts ids %v failed: %w", hg.Hosts, err)
	}
	return hosts, nil
}

// HostInUse find host is used by hostgroup
func HostInUse(tx *gorm.DB, id string) (bool, error) {
	var count int64
	err := tx.Where("hosts LIKE ?", "%"+id+"%").Find(&HostGroup{}).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("hosts like failed: %w", err)
	}
	return count > 0, nil
}

// HostInUse find host is used by hostgroup
func HostGroupInUse(tx *gorm.DB, id string) (bool, error) {
	var count int64
	err := tx.Where("hostgroup_id = ?", id).Find(&Task{}).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("hosts like failed: %w", err)
	}
	return count > 0, nil
}
